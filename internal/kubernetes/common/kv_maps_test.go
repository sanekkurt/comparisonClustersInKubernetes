package common

import (
	"context"
	"errors"
	"fmt"
	"k8s-cluster-comparator/internal/kubernetes/diff"
	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func initCtx() context.Context {
	var (
		ctx = context.Background()
	)
	err := logging.ConfigureForTests()
	if err != nil {
		fmt.Println("[ERROR] ", err.Error())
	}

	return ctx
}
func newCtxWithCleanStorage(ctx context.Context) context.Context {
	diffs := diff.NewDiffsStorage(ctx)
	ctx = diff.WithDiffStorage(ctx, diffs)

	batch := diffs.NewLazyBatch(metav1.TypeMeta{Kind: "", APIVersion: ""}, metav1.ObjectMeta{})
	ctx = diff.WithDiffBatch(ctx, batch)

	return ctx
}

func initKVMapsForTest1() (types.KVMap, types.KVMap, map[string]struct{}, bool) {
	kvm1 := make(types.KVMap)
	kvm2 := make(types.KVMap)

	kvm1["key1"] = "value1"
	kvm1["key2"] = "value2"
	kvm1["key3"] = "value3"
	kvm1["key4"] = "value4"

	kvm2["key1"] = "value1"
	kvm2["key2"] = "value2"
	kvm2["key3"] = "value3"

	skipKeys := make(map[string]struct{})
	skipKeys["key3"] = struct{}{}

	return kvm1, kvm2, skipKeys, true
}

func initKVMapsForTest2() (types.KVMap, types.KVMap, map[string]struct{}, bool) {
	kvm1 := make(types.KVMap)
	kvm2 := make(types.KVMap)

	kvm1["key1"] = "value1"
	kvm1["key2"] = "value2"
	kvm1["key3"] = "value3"

	kvm2["key1"] = "value1"
	kvm2["key2"] = "value2"
	kvm2["key3"] = "value3"
	kvm2["key4"] = "value4"

	skipKeys := make(map[string]struct{})
	skipKeys["key3"] = struct{}{}

	return kvm1, kvm2, skipKeys, true
}

func initKVMapsForTest3() (types.KVMap, types.KVMap, map[string]struct{}, bool) {
	kvm1 := make(types.KVMap)
	kvm2 := make(types.KVMap)

	kvm1["key1"] = "value1"
	kvm1["key2"] = "diffValue2"
	kvm1["key3"] = "value3"

	kvm2["key1"] = "value1"
	kvm2["key2"] = "value2"
	kvm2["key3"] = "diffValue"

	skipKeys := make(map[string]struct{})
	skipKeys["key3"] = struct{}{}

	return kvm1, kvm2, skipKeys, true
}

func initKVMapsForTest4() (types.KVMap, types.KVMap, map[string]struct{}, bool) {
	kvm1 := make(types.KVMap)
	kvm2 := make(types.KVMap)

	kvm1["key1"] = "value1"
	kvm1["key2"] = "diffValue2"
	kvm1["key3"] = "value3"

	kvm2["key1"] = "value1"
	kvm2["key2"] = "value2"
	kvm2["key3"] = "diffValue"

	skipKeys := make(map[string]struct{})
	skipKeys["key3"] = struct{}{}

	return kvm1, kvm2, skipKeys, false
}

func initKVBytesMapsForTest1() (map[string][]byte, map[string][]byte, map[string]struct{}, bool) {
	kvm1 := make(map[string][]byte)
	kvm2 := make(map[string][]byte)

	kvm1["key1"] = []byte("value1")
	kvm1["key2"] = []byte("value2")
	kvm1["key3"] = []byte("value3")
	kvm1["key4"] = []byte("value4")

	kvm2["key1"] = []byte("value1")
	kvm2["key2"] = []byte("value2")
	kvm2["key3"] = []byte("value3")

	skipKeys := make(map[string]struct{})
	skipKeys["key3"] = struct{}{}

	return kvm1, kvm2, skipKeys, true
}

func initKVBytesMapsForTest2() (map[string][]byte, map[string][]byte, map[string]struct{}, bool) {
	kvm1 := make(map[string][]byte)
	kvm2 := make(map[string][]byte)

	kvm1["key1"] = []byte("value1")
	kvm1["key2"] = []byte("value2")
	kvm1["key3"] = []byte("value3")

	kvm2["key1"] = []byte("value1")
	kvm2["key2"] = []byte("value2")
	kvm2["key3"] = []byte("value3")
	kvm2["key4"] = []byte("value4")

	skipKeys := make(map[string]struct{})
	skipKeys["key3"] = struct{}{}

	return kvm1, kvm2, skipKeys, true
}

func initKVBytesMapsForTest3() (map[string][]byte, map[string][]byte, map[string]struct{}, bool) {
	kvm1 := make(map[string][]byte)
	kvm2 := make(map[string][]byte)

	kvm1["key1"] = []byte("value1")
	kvm1["key2"] = []byte("diffValue2")
	kvm1["key3"] = []byte("value3")

	kvm2["key1"] = []byte("value1")
	kvm2["key2"] = []byte("value2")
	kvm2["key3"] = []byte("diffValue")

	skipKeys := make(map[string]struct{})
	skipKeys["key3"] = struct{}{}

	return kvm1, kvm2, skipKeys, true
}

func initKVBytesMapsForTest4() (map[string][]byte, map[string][]byte, map[string]struct{}, bool) {
	kvm1 := make(map[string][]byte)
	kvm2 := make(map[string][]byte)

	kvm1["key1"] = []byte("value1")
	kvm1["key2"] = []byte("diffValue2")
	kvm1["key3"] = []byte("value3")

	kvm2["key1"] = []byte("value1")
	kvm2["key2"] = []byte("value2")
	kvm2["key3"] = []byte("diffValue")

	skipKeys := make(map[string]struct{})
	skipKeys["key3"] = struct{}{}

	return kvm1, kvm2, skipKeys, false
}

func TestAreKVMapsEqual(t *testing.T) {
	cleanCtx := initCtx()

	ctx := newCtxWithCleanStorage(cleanCtx)
	hostPath1, hostPath2, skipKeys, dumpValues := initKVMapsForTest1()

	equal := AreKVMapsEqual(ctx, hostPath1, hostPath2, skipKeys, dumpValues)

	diffStorage := diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch := diff.BatchFromContext(ctx)
	diffs := batch.GetDiffs()

	if diffs != nil && !equal {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorKeyDoesNotExistInMap2) {
				t.Errorf("Error expected: '%s: 'key4''. But it was returned: %s", ErrorKeyDoesNotExistInMap2.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else if equal {
		t.Errorf("function returned true but should have returned false")
	} else {
		t.Errorf("Error expected: '%s: 'key4''. But the function found no errors", ErrorKeyDoesNotExistInMap2.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	hostPath1, hostPath2, skipKeys, dumpValues = initKVMapsForTest2()

	equal = AreKVMapsEqual(ctx, hostPath1, hostPath2, skipKeys, dumpValues)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil && !equal {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorExtraKeysFoundInMap2) {
				t.Errorf("Error expected: '%s: '1''. But it was returned: %s", ErrorExtraKeysFoundInMap2.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else if equal {
		t.Errorf("function returned true but should have returned false")
	} else {
		t.Errorf("Error expected: '%s: '1''. But the function found no errors", ErrorExtraKeysFoundInMap2.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	hostPath1, hostPath2, skipKeys, dumpValues = initKVMapsForTest3()

	equal = AreKVMapsEqual(ctx, hostPath1, hostPath2, skipKeys, dumpValues)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil && !equal {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorKeyValueDoNotMatchInMap) {
				t.Errorf("Error expected: '%s. key2: 'diffValue2' vs 'value2''. But it was returned: %s", ErrorKeyValueDoNotMatchInMap.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else if equal {
		t.Errorf("function returned true but should have returned false")
	} else {
		t.Errorf("Error expected: '%s. key2: 'diffValue2' vs 'value2''. But the function found no errors", ErrorKeyValueDoNotMatchInMap.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	hostPath1, hostPath2, skipKeys, dumpValues = initKVMapsForTest4()

	equal = AreKVMapsEqual(ctx, hostPath1, hostPath2, skipKeys, dumpValues)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil && !equal {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorKeyValueDoNotMatchInMap) {
				t.Errorf("Error expected: '%s. Key: 'key2''. But it was returned: %s", ErrorKeyValueDoNotMatchInMap.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else if equal {
		t.Errorf("function returned true but should have returned false")
	} else {
		t.Errorf("Error expected: '%s. Key: 'key2''. But the function found no errors", ErrorKeyValueDoNotMatchInMap.Error())
	}
}

func TestAreKVBytesMapsEqual(t *testing.T) {
	cleanCtx := initCtx()

	ctx := newCtxWithCleanStorage(cleanCtx)
	hostPath1, hostPath2, skipKeys, dumpValues := initKVBytesMapsForTest1()

	equal := AreKVBytesMapsEqual(ctx, hostPath1, hostPath2, skipKeys, dumpValues)

	diffStorage := diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch := diff.BatchFromContext(ctx)
	diffs := batch.GetDiffs()

	if diffs != nil && !equal {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorKeyDoesNotExistInMap2) {
				t.Errorf("Error expected: '%s: 'key4''. But it was returned: %s", ErrorKeyDoesNotExistInMap2.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else if equal {
		t.Errorf("function returned true but should have returned false")
	} else {
		t.Errorf("Error expected: '%s: 'key4''. But the function found no errors", ErrorKeyDoesNotExistInMap2.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	hostPath1, hostPath2, skipKeys, dumpValues = initKVBytesMapsForTest2()

	equal = AreKVBytesMapsEqual(ctx, hostPath1, hostPath2, skipKeys, dumpValues)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil && !equal {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorExtraKeysFoundInMap2) {
				t.Errorf("Error expected: '%s: '1''. But it was returned: %s", ErrorExtraKeysFoundInMap2.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else if equal {
		t.Errorf("function returned true but should have returned false")
	} else {
		t.Errorf("Error expected: '%s: '1''. But the function found no errors", ErrorExtraKeysFoundInMap2.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	hostPath1, hostPath2, skipKeys, dumpValues = initKVBytesMapsForTest3()

	equal = AreKVBytesMapsEqual(ctx, hostPath1, hostPath2, skipKeys, dumpValues)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil && !equal {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorKeyValueDoNotMatchInMap) {
				t.Errorf("Error expected: '%s. key2: 'diffValue2' vs 'value2''. But it was returned: %s", ErrorKeyValueDoNotMatchInMap.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else if equal {
		t.Errorf("function returned true but should have returned false")
	} else {
		t.Errorf("Error expected: '%s. key2: 'diffValue2' vs 'value2''. But the function found no errors", ErrorKeyValueDoNotMatchInMap.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	hostPath1, hostPath2, skipKeys, dumpValues = initKVBytesMapsForTest4()

	equal = AreKVBytesMapsEqual(ctx, hostPath1, hostPath2, skipKeys, dumpValues)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil && !equal {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorKeyValueDoNotMatchInMap) {
				t.Errorf("Error expected: '%s. Key: 'key2''. But it was returned: %s", ErrorKeyValueDoNotMatchInMap.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else if equal {
		t.Errorf("function returned true but should have returned false")
	} else {
		t.Errorf("Error expected: '%s. Key: 'key2''. But the function found no errors", ErrorKeyValueDoNotMatchInMap.Error())
	}
}
