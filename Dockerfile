FROM ubuntu:latest as build
FROM index.docker.io/library/python:3.9-slim@sha256:27211e8bbfd2c91ac9adbde0565b9ac18234bfcde8ef0e6a3404fd404f26ea13
FROM node:latest
FROM index.docker.io/library/nginx:alpine@sha256:208b70eefac13ee9be00e486f79c695b15cef861c680527171a27d253d834be9
