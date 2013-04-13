package main

import (
	"errors"
)

func valid_chan(id_chan, id_user int, need_connected bool) error {
	var err error

	if user, ok := all_users[id_user]; ok {
		if chann, ok := user.Buffers[id_chan]; ok && chann.connected == need_connected {
			if serv, ok := user.Buffers[chann.id_serv]; ok && serv.connected == need_connected {

			} else {
				err = errors.New("Chan : invalid server")
			}
		} else {
			err = errors.New("Chan : invalid chan")
		}
	} else {
		err = errors.New("Chan : invalid user")
	}
	return err
}

func valid_serv(id_serv, id_user int, need_connected bool) error {
	var err error

	if user, ok := all_users[id_user]; ok {
		if serv, ok := user.Buffers[id_serv]; ok && serv.connected == need_connected {

		} else {
			err = errors.New("Serv : invalid server")
		}
	} else {
		err = errors.New("Serv : invalid user")
	}
	return err
}
