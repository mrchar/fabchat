package main

import "github.com/pkg/errors"

func checkParamsLength(params []string, length int) error {
	if len(params) != length {
		return errors.Errorf("需要%d个参数，得到%d个参数", length, len(params))
	}
	return nil
}
