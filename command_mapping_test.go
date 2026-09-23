package main

import (
	"hash/crc32"
	"reflect"
	"testing"
)

func crc16(name string) int {
	return int(uint16(crc32.ChecksumIEEE([]byte(name)) & 0xFFFF))
}

func copyMapping(m map[string]int) map[string]int {
	c := make(map[string]int, len(m))
	for k, v := range m {
		c[k] = v
	}
	return c
}

// TestResolveCommandMapping_SelfHeal 历史人工值在crc空闲时自愈回归crc值,且幂等
func TestResolveCommandMapping_SelfHeal(t *testing.T) {
	names := []string{"PlayerReconnectGameReq", "PlayerReconnectGameRes", "LoginReq"}
	old := map[string]int{
		"PlayerReconnectGameReq":  4, // 历史碰撞分配的人工值
		"PlayerReconnectGameRes":  3,
		"LoginReq":                crc16("LoginReq"),
	}
	got := resolveCommandMapping(names, old)
	for _, n := range names {
		if got[n] != crc16(n) {
			t.Errorf("%v: expect self-heal to crc %v, got %v", n, crc16(n), got[n])
		}
	}
	// 幂等:同样输入再来一次,结果完全一致
	again := resolveCommandMapping(names, copyMapping(got))
	if !reflect.DeepEqual(got, again) {
		t.Errorf("not idempotent:\nfirst=%v\nsecond=%v", got, again)
	}
}

// TestResolveCommandMapping_KeepManualValueWhenCrcOccupied crc被他人占用时人工值保持稳定(不再漂移);
// 占用者移除后自愈回归crc值(至多两轮收敛,收敛后幂等)
func TestResolveCommandMapping_KeepManualValueWhenCrcOccupied(t *testing.T) {
	a, b := "MsgA", "MsgB"
	old := map[string]int{
		a: 7,        // A的人工值
		b: crc16(a), // B占用着A的crc值(模拟A的crc被他人占用的碰撞场景)
	}
	got1 := resolveCommandMapping([]string{a, b}, old)
	// A的crc被B占用,A必须保持人工值7,不能被重新分配(旧实现的漂移点)
	if got1[a] != 7 {
		t.Errorf("MsgA: expect keep manual value 7, got %v", got1[a])
	}
	if got1[b] != crc16(b) {
		t.Errorf("MsgB: expect crc value %v, got %v", crc16(b), got1[b])
	}
	// 第二轮:伪造占用已被纠正,全部回归crc值
	got2 := resolveCommandMapping([]string{a, b}, copyMapping(got1))
	if got2[a] != crc16(a) || got2[b] != crc16(b) {
		t.Errorf("expect converged to crc values, got %v", got2)
	}
	// 收敛后幂等
	got3 := resolveCommandMapping([]string{a, b}, copyMapping(got2))
	if !reflect.DeepEqual(got2, got3) {
		t.Errorf("not idempotent after converge:\n%v\n%v", got2, got3)
	}
	// 删除消息场景:释放的旧消息不残留
	removed := resolveCommandMapping([]string{a}, copyMapping(got3))
	if len(removed) != 1 {
		t.Errorf("expect removed stale message, got %v", removed)
	}
}
