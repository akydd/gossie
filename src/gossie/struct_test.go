package gossie

import (
    "testing"
    "reflect"
    "os"
)

/*

todo:

    basically everything. real unit testing for all struct funcs

*/

type errNoMeta struct {
    a int
}
type errNoMetaKeyColVal struct {
    a int `cf:"cfname"`
}
type errNoMetaColVal struct {
    a int `cf:"cfname" key:"a"`
}
type errNoMetaVal struct {
    a   int `cf:"cfname" key:"a" col:"b"`
    b   int
}
type errInvKey struct {
    a   int `cf:"cfname" key:"z" col:"b" val:"c"`
    b   int
    c   int
}
type errInvCol struct {
    a   int `cf:"cfname" key:"a" col:"z" val:"c"`
    b   int
    c   int
}
type errInvVal struct {
    a   int `cf:"cfname" key:"a" col:"b" val:"z"`
    b   int
    c   int
}
type noErrA struct {
    a   int `cf:"cfname" key:"a" col:"b" val:"c"`
    b   int
    c   int
}
type noErrB struct {
    a   int `cf:"cfname" key:"a" col:"*name" val:"*value"`
    b   int
    c   int
}
type noErrC struct {
    a   int `cf:"cfname" key:"a" col:"b,*name" val:"*value"`
    b   int
    c   int
}
type noErrD struct {
    a   int `cf:"cfname" key:"a" col:"b" val:"c"`
    b   []int
    c   []int
}
type noErrE struct {
    a   int `cf:"cfname" key:"a" col:"b,c" val:"d"`
    b   int
    c   []int
    d   []int
}
type everythingComp struct {
    Key      string `cf:"cfname" key:"Key" col:"FBytes,FBool,FInt8,FInt16,FInt32,FInt,FInt64,FFloat32,FFloat64,FString,FUUID,*name" val:"*value"`
    FBytes   []byte
    FBool    bool
    FInt8    int8
    FInt16   int16
    FInt32   int32
    FInt     int
    FInt64   int64
    FFloat32 float32
    FFloat64 float64
    FString  string
    FUUID    UUID
    Val      string
}

func buildMappingFromPtr(instance interface{}) (*structMapping, os.Error) {
    valuePtr := reflect.ValueOf(instance)
    value := reflect.Indirect(valuePtr)
    typ := value.Type()
    return newStructMapping(typ)
}

func structMapMustError(t *testing.T, instance interface{}) {
    _, err := buildMappingFromPtr(instance)
    if err == nil {
        t.Error("Expected error calling newStructMapping, got none")
    }
}

func checkMapping(t *testing.T, expected, actual interface{}, name string) {
    if !reflect.DeepEqual(expected, actual) {
        t.Error("Mapping for struct sample", name, "does not match expected output")
    }
}

func TestStructMapping(t *testing.T) {
    structMapMustError(t, &errNoMeta{})
    structMapMustError(t, &errNoMetaKeyColVal{})
    structMapMustError(t, &errNoMetaColVal{})
    structMapMustError(t, &errNoMetaVal{})
    structMapMustError(t, &errInvKey{})
    structMapMustError(t, &errInvCol{})
    structMapMustError(t, &errInvVal{})

    mapA, _ := buildMappingFromPtr(&noErrA{1, 2, 3})
    goodA := &structMapping{
        cf:  "cfname",
        key: &fieldMapping{fieldKind: baseTypeField, position: 0, name: "a", cassandraType: LongType},
        columns: []*fieldMapping{
            &fieldMapping{fieldKind: baseTypeField, position: 1, name: "b", cassandraType: LongType},
        },
        value:             &fieldMapping{fieldKind: baseTypeField, position: 2, name: "c", cassandraType: LongType},
        others:            nil,
        isCompositeColumn: false,
    }
    checkMapping(t, goodA, mapA, "mapA")

    mapB, _ := buildMappingFromPtr(&noErrB{1, 2, 3})
    goodB := &structMapping{
        cf:  "cfname",
        key: &fieldMapping{fieldKind: baseTypeField, position: 0, name: "a", cassandraType: LongType},
        columns: []*fieldMapping{
            &fieldMapping{fieldKind: starNameField, position: 0, name: "", cassandraType: 0},
        },
        value: &fieldMapping{fieldKind: starValueField, position: 0, name: "", cassandraType: 0},
        others: map[string]*fieldMapping{
            "b": &fieldMapping{fieldKind: baseTypeField, position: 1, name: "b", cassandraType: LongType},
            "c": &fieldMapping{fieldKind: baseTypeField, position: 2, name: "c", cassandraType: LongType},
        },
        isCompositeColumn: false,
    }
    checkMapping(t, goodB, mapB, "mapB")

    mapC, _ := buildMappingFromPtr(&noErrC{1, 2, 3})
    goodC := &structMapping{
        cf:  "cfname",
        key: &fieldMapping{fieldKind: baseTypeField, position: 0, name: "a", cassandraType: LongType},
        columns: []*fieldMapping{
            &fieldMapping{fieldKind: baseTypeField, position: 1, name: "b", cassandraType: LongType},
            &fieldMapping{fieldKind: starNameField, position: 0, name: "", cassandraType: 0},
        },
        value: &fieldMapping{fieldKind: starValueField, position: 0, name: "", cassandraType: 0},
        others: map[string]*fieldMapping{
            "c": &fieldMapping{fieldKind: baseTypeField, position: 2, name: "c", cassandraType: LongType},
        },
        isCompositeColumn: true,
    }
    checkMapping(t, goodC, mapC, "mapC")

    mapD, _ := buildMappingFromPtr(&noErrD{1, []int{2, 3}, []int{4, 5}})
    goodD := &structMapping{
        cf:  "cfname",
        key: &fieldMapping{fieldKind: baseTypeField, position: 0, name: "a", cassandraType: LongType},
        columns: []*fieldMapping{
            &fieldMapping{fieldKind: baseTypeSliceField, position: 1, name: "b", cassandraType: LongType},
        },
        value:             &fieldMapping{fieldKind: baseTypeSliceField, position: 2, name: "c", cassandraType: LongType},
        others:            nil,
        isCompositeColumn: false,
    }
    checkMapping(t, goodD, mapD, "mapD")

    mapE, _ := buildMappingFromPtr(&noErrE{1, 2, []int{3, 4}, []int{5, 6}})
    goodE := &structMapping{
        cf:  "cfname",
        key: &fieldMapping{fieldKind: baseTypeField, position: 0, name: "a", cassandraType: LongType},
        columns: []*fieldMapping{
            &fieldMapping{fieldKind: baseTypeField, position: 1, name: "b", cassandraType: LongType},
            &fieldMapping{fieldKind: baseTypeSliceField, position: 2, name: "c", cassandraType: LongType},
        },
        value:             &fieldMapping{fieldKind: baseTypeSliceField, position: 3, name: "d", cassandraType: LongType},
        others:            nil,
        isCompositeColumn: true,
    }
    checkMapping(t, goodE, mapE, "mapE")

    eComp, _ := buildMappingFromPtr(&everythingComp{"a", []byte{1, 2}, true, 3, 4, 5, 6, 7, 8.0, 9.0, "b",
        [16]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}, "c"})
    goodEComp := &structMapping{
        cf:  "cfname",
        key: &fieldMapping{fieldKind: baseTypeField, position: 0, name: "Key", cassandraType: UTF8Type},
        columns: []*fieldMapping{
            &fieldMapping{fieldKind: baseTypeField, position: 1, name: "FBytes", cassandraType: BytesType},
            &fieldMapping{fieldKind: baseTypeField, position: 2, name: "FBool", cassandraType: BooleanType},
            &fieldMapping{fieldKind: baseTypeField, position: 3, name: "FInt8", cassandraType: LongType},
            &fieldMapping{fieldKind: baseTypeField, position: 4, name: "FInt16", cassandraType: LongType},
            &fieldMapping{fieldKind: baseTypeField, position: 5, name: "FInt32", cassandraType: LongType},
            &fieldMapping{fieldKind: baseTypeField, position: 6, name: "FInt", cassandraType: LongType},
            &fieldMapping{fieldKind: baseTypeField, position: 7, name: "FInt64", cassandraType: LongType},
            &fieldMapping{fieldKind: baseTypeField, position: 8, name: "FFloat32", cassandraType: FloatType},
            &fieldMapping{fieldKind: baseTypeField, position: 9, name: "FFloat64", cassandraType: DoubleType},
            &fieldMapping{fieldKind: baseTypeField, position: 10, name: "FString", cassandraType: UTF8Type},
            &fieldMapping{fieldKind: baseTypeField, position: 11, name: "FUUID", cassandraType: UUIDType},
            &fieldMapping{fieldKind: starNameField, position: 0, name: "", cassandraType: 0},
        },
        value: &fieldMapping{fieldKind: starValueField, position: 0, name: "", cassandraType: 0},
        others: map[string]*fieldMapping{
            "Val": &fieldMapping{fieldKind: baseTypeField, position: 12, name: "Val", cassandraType: UTF8Type},
        },
        isCompositeColumn: true,
    }
    t.Log(goodEComp.key)
    t.Log(eComp.key)
    checkMapping(t, goodEComp, eComp, "eComp")

}

func TestMap(t *testing.T) {
    ec := &everythingComp{"a", []byte{1, 2}, true, 3, 4, 5, 6, 7, 8.0, 9.0, "b",
        [16]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}, "c"}
    row, err := Map(ec)
    if err != nil {
        t.Fatal("Unexpected error in test map:", err)
    }
    if len(row.Columns) != 1 {
        t.Error("Expected number of columns is 1, got ", len(row.Columns))
    }
    name := []byte{0, 2, 1, 2, 0, 0, 1, 1, 0, 0, 8, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 8, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 8, 0, 0, 0,
        0, 0, 0, 0, 5, 0, 0, 8, 0, 0, 0, 0, 0, 0, 0, 6, 0, 0, 8, 0, 0, 0, 0, 0, 0, 0, 7, 0, 0, 4, 65, 0, 0, 0, 0, 0,
        8, 64, 34, 0, 0, 0, 0, 0, 0, 0, 0, 1, 98, 0, 0, 16, 0, 17, 34, 51, 68, 85, 102, 119, 136, 153, 170, 187, 204,
        221, 238, 255, 0, 0, 3, 86, 97, 108, 0}
    value := []byte{99}
    if !reflect.DeepEqual(name, row.Columns[0].Name) {
        t.Error("Invalid composite column name in for test row")
    }
    if !reflect.DeepEqual(value, row.Columns[0].Value) {
        t.Error("Invalid value in test row")
    }
}
