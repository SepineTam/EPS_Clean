#!/usr/bin/python3
# -*- coding: utf-8 -*-
#
# Copyright (C) 2025 - Present Sepine Tam, Inc. All Rights Reserved
#
# @Author : Sepine Tam
# @Email  : sepinetam@gmail.com
# @File   : setup.py

from setuptools import setup, find_packages

setup(
    name="epsclean",
    version="0.1.0",
    packages=find_packages(),
    install_requires=[
        "pandas",
    ],
    entry_points={
        'console_scripts': [
            'epsclean=epsclean.main:main',
        ],
    },
    author="Song Tan",
    author_email="sepinetam@gmail.com",
    description="A tool to clean CSV files by extracting the last three rows and converting to UTF-8",
    keywords="csv, cleaning, encoding",
    python_requires='>=3.6',
)