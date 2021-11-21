#!/usr/bin/env python
# coding=utf-8
# Author:Wa-llee

from setuptools import setup,find_packages
import os
setup(
    name='chinese2digits',
    version= '6.0.1',
    keywords = ("pip", "NLP","chinese","digits_number","number_transform"),
    description=(
        """最好的汉字数字(中文数字)-阿拉伯数字转换工具。包含"点二八"，"负百分之四十"等众多汉语表达方法。NLP，机器人工程必备！ The Best Tool of Chinese Number to Digits"""
    ),
    long_description=open(
        os.path.join(
            os.path.dirname(__file__),
            'README.rst'
        )
    ).read(),
    author='Wa-llee',
    author_email='xrli_office@foxmail.com',
    license='Apache License 2.0',
    packages= find_packages(),
    platforms=["all"],
    url='https://github.com/Wall-ee/chinese2digits'
)

