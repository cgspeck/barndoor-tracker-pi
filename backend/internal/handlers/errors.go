package handlers

import "fmt"

type UnrecognisedPathError struct {
	Path string
}

func (u UnrecognisedPathError) Error() string {
	return fmt.Sprintf("Unrecognised path: %v", u.Path)
}

type BadRequestError struct{}

func (_ BadRequestError) Error() string {
	return "Bad Request"
}

type CouldNotCastToBoolError struct {
	val interface{}
}

func (e CouldNotCastToBoolError) Error() string {
	return fmt.Sprintf("Could not cast %v to bool", e.val)
}

type CouldNotCastToFloat64Error struct {
	val interface{}
}

func (e CouldNotCastToFloat64Error) Error() string {
	return fmt.Sprintf("Could not cast %v to float64", e.val)
}

type CouldNotCastToIntError struct {
	val interface{}
}

func (e CouldNotCastToIntError) Error() string {
	return fmt.Sprintf("Could not cast %v to int", e.val)
}

type KeyNotFoundError struct {
	k       string
	hashMap interface{}
}

func (e KeyNotFoundError) Error() string {
	return fmt.Sprintf("Could not find %q in map %+v", e.k, e.hashMap)
}
