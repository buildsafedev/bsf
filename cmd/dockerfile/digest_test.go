package dockerfile

import (
	"bytes"
	"testing"
)

func TestReadByte(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected []string
	}{
		{
			name: "Test 1",
			input: []byte(`
				FROM ubuntu:latest
				FROM nginx:1.19.10
				# Comment line for Testing
				FROM golang:1.16.5-alpine3.13
			`),
			expected: []string{"ubuntu:latest", "nginx:1.19.10", "golang:1.16.5-alpine3.13"},
		},
		{
			name:     "Test 2",
			input:    []byte(""),
			expected: []string{},
		},
		{
			name: "Test 3",
			input: []byte(`
				# Testing ....
				RUN apt-get update
			`),
			expected: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lines, err := readByte(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !equalSlice(lines, tc.expected) {
				t.Errorf("got %v; expected %v", lines, tc.expected)
			}
		})
	}
}

func equalSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestUpdatedockerfileWithDigests(t *testing.T) {
	testCases := []struct {
		name         string
		inputData    []byte
		digestMap    map[string]string
		expectedData []byte
	}{
		{
			name: "Test 1",
			inputData: []byte(`
				FROM ubuntu:20.04
			`),
			digestMap: map[string]string{
				"ubuntu:20.04": "sha256:0b897358ff6624825fb50d20ffb605ab0eaea77ced0adb8c6a4b756513dec6fc",
			},
			expectedData: []byte(`
				FROM ubuntu@sha256:0b897358ff6624825fb50d20ffb605ab0eaea77ced0adb8c6a4b756513dec6fc
			`),
		},
		{
			name: "Test 2",
			inputData: []byte(`
				FROM ubuntu:20.04
				FROM python:3.9-slim
				FROM node:14
				FROM node:latest AS build
			`),
			digestMap: map[string]string{
				"ubuntu:20.04": "sha256:0b897358ff6624825fb50d20ffb605ab0eaea77ced0adb8c6a4b756513dec6fc",
				"node:14":      "sha256:a158d3b9b4e3fa813fa6c8c590b8f0a860e015ad4e59bbce5744d2f6fd8461aa",
				"node:latest":  "sha256:c8a559f733bf1f9b3c1d05b97d9a9c7e5d3647c99abedaf5cdd3b54c9cbb8eff",
			},
			expectedData: []byte(`
				FROM ubuntu@sha256:0b897358ff6624825fb50d20ffb605ab0eaea77ced0adb8c6a4b756513dec6fc
				FROM python:3.9-slim
				FROM node@sha256:a158d3b9b4e3fa813fa6c8c590b8f0a860e015ad4e59bbce5744d2f6fd8461aa
				FROM node@sha256:c8a559f733bf1f9b3c1d05b97d9a9c7e5d3647c99abedaf5cdd3b54c9cbb8eff AS build
			`),
		},
		{
			name: "Test 3",
			inputData: []byte(`
				FROM busybox:latest
				RUN apt-get update
			`),
			digestMap: map[string]string{
				"ubuntu:20.04": "sha256:0b897358ff6624825fb50d20ffb605ab0eaea77ced0adb8c6a4b756513dec6fc",
			},
			expectedData: []byte(`
				FROM busybox:latest
				RUN apt-get update
			`),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := updateDockerfileWithDigests(tc.inputData, tc.digestMap)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !bytes.Equal(result, tc.expectedData) {
				t.Errorf("got %s; expected %s", result, tc.expectedData)
			}
		})
	}
}
