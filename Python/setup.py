#!/usr/bin/env python
# coding=utf-8
# Author:Wa-llee

from setuptools import setup,find_packages

setup(
    name='chinese2digits',
    version= '1.0.1',
    description=(
        """最好的汉字数字(中文数字)-阿拉伯数字转换工具。包含"点二八"，"负百分之四十"等众多汉语表达方法。NLP，机器人工程必备！ The Best Tool of Chinese Number to Digits"""
    ),
    author='Wa-llee',
    author_email='xrli_office@foxmail.com',
    license='Apache License 2.0',
    packages= find_packages(),
    platforms=["all"],
    url='https://github.com/Wall-ee/chinese2digits'
)