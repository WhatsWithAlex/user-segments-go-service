package koanfext

import (
	"github.com/knadh/koanf/v2"
	"github.com/mitchellh/mapstructure"
)

type KoanfExt struct {
	*koanf.Koanf
}

func New(delim string) *KoanfExt {
	return &KoanfExt{
		koanf.New(delim),
	}
}

func (ke *KoanfExt) StrictUnmarshal(path string, o interface{}) error {
	c := &mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.TextUnmarshallerHookFunc()),
		Metadata:         nil,
		Result:           o,
		WeaklyTypedInput: true,
		ErrorUnset:       true,
		ErrorUnused:      true,
	}
	return ke.UnmarshalWithConf(path, o, koanf.UnmarshalConf{DecoderConfig: c})
}

func (ke *KoanfExt) StrictUnmarshalWithHook(path string, o interface{}, decodeHook mapstructure.DecodeHookFunc) error {
	c := &mapstructure.DecoderConfig{
		DecodeHook:       decodeHook,
		Metadata:         nil,
		Result:           o,
		WeaklyTypedInput: true,
		ErrorUnset:       true,
		ErrorUnused:      true,
	}
	return ke.UnmarshalWithConf(path, o, koanf.UnmarshalConf{DecoderConfig: c})
}
