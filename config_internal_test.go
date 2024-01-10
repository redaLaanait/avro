package avro

import (
	"testing"

	"github.com/modern-go/reflect2"
	"github.com/stretchr/testify/assert"
)

func TestConfig_Freeze(t *testing.T) {
	api := Config{
		TagKey:      "test",
		BlockLength: 2,
	}.Freeze()
	cfg := api.(*frozenConfig)

	assert.Equal(t, "test", cfg.getTagKey())
	assert.Equal(t, 2, cfg.getBlockLength())
}

func TestConfig_ReusesDecoders(t *testing.T) {
	type testObj struct {
		A int64 `avro:"a"`
	}

	api := Config{
		TagKey:      "test",
		BlockLength: 2,
	}.Freeze()
	cfg := api.(*frozenConfig)

	schema := MustParse(`{
	"type": "record",
	"name": "test",
	"fields" : [
		{"name": "a", "type": "long"}
	]
}`)
	typ := reflect2.TypeOfPtr(&testObj{})

	dec1 := cfg.DecoderOf(schema, typ)
	dec2 := cfg.DecoderOf(schema, typ)

	assert.Same(t, dec1, dec2)
}

func TestConfig_ReusesDecoders_WithRecordFieldActions(t *testing.T) {
	type testObj struct {
		A int64  `avro:"a"`
		B string `avro:"b"`
	}
	sch := `{
		"type": "record",
		"name": "test",
		"fields" : [
			{"name": "a", "type": "long"},
			{"name": "a", "type": "string", "default": "foo"}
		]
	}`
	typ := reflect2.TypeOfPtr(&testObj{})

	t.Run("set default", func(t *testing.T) {
		api := Config{
			TagKey:      "test",
			BlockLength: 2,
		}.Freeze()
		cfg := api.(*frozenConfig)

		schema1 := MustParse(sch)
		schema2 := MustParse(sch)
		schema2.(*RecordSchema).Fields()[1].action = FieldSetDefault

		dec1 := cfg.DecoderOf(schema1, typ)
		dec2 := cfg.DecoderOf(schema2, typ)

		assert.NotSame(t, dec1, dec2)
	})

	t.Run("ignore", func(t *testing.T) {
		api := Config{
			TagKey:      "test",
			BlockLength: 2,
		}.Freeze()
		cfg := api.(*frozenConfig)

		schema1 := MustParse(sch)
		schema1.(*RecordSchema).Fields()[0].action = FieldIgnore
		schema2 := MustParse(sch)

		dec1 := cfg.DecoderOf(schema1, typ)
		dec2 := cfg.DecoderOf(schema2, typ)

		assert.NotSame(t, dec1, dec2)
	})

}

func TestConfig_DisableCache_DoesNotReuseDecoders(t *testing.T) {
	type testObj struct {
		A int64 `avro:"a"`
	}

	api := Config{
		TagKey:         "test",
		BlockLength:    2,
		DisableCaching: true,
	}.Freeze()
	cfg := api.(*frozenConfig)

	schema := MustParse(`{
	"type": "record",
	"name": "test",
	"fields" : [
		{"name": "a", "type": "long"}
	]
}`)
	typ := reflect2.TypeOfPtr(&testObj{})

	dec1 := cfg.DecoderOf(schema, typ)
	dec2 := cfg.DecoderOf(schema, typ)

	assert.NotSame(t, dec1, dec2)
}

func TestConfig_ReusesEncoders(t *testing.T) {
	type testObj struct {
		A int64 `avro:"a"`
	}

	api := Config{
		TagKey:      "test",
		BlockLength: 2,
	}.Freeze()
	cfg := api.(*frozenConfig)

	schema := MustParse(`{
	"type": "record",
	"name": "test",
	"fields" : [
		{"name": "a", "type": "long"}
	]
}`)
	typ := reflect2.TypeOfPtr(testObj{})

	enc1 := cfg.EncoderOf(schema, typ)
	enc2 := cfg.EncoderOf(schema, typ)

	assert.Same(t, enc1, enc2)
}

func TestConfig_DisableCache_DoesNotReuseEncoders(t *testing.T) {
	type testObj struct {
		A int64 `avro:"a"`
	}

	api := Config{
		TagKey:         "test",
		BlockLength:    2,
		DisableCaching: true,
	}.Freeze()
	cfg := api.(*frozenConfig)

	schema := MustParse(`{
	"type": "record",
	"name": "test",
	"fields" : [
		{"name": "a", "type": "long"}
	]
}`)
	typ := reflect2.TypeOfPtr(testObj{})

	enc1 := cfg.EncoderOf(schema, typ)
	enc2 := cfg.EncoderOf(schema, typ)

	assert.NotSame(t, enc1, enc2)
}
