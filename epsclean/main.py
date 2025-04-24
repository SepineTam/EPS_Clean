#!/usr/bin/python3
# -*- coding: utf-8 -*-
#
# Copyright (C) 2025 - Present Sepine Tam, Inc. All Rights Reserved
#
# @Author : Sepine Tam
# @Email  : sepinetam@gmail.com
# @File   : main.py

import pandas as pd
import argparse
import os
import sys


def process_csv(input_file, output_file=None, encoding='gb2312'):
    """
    Read CSV file, remove the last three rows, and save with UTF-8 encoding

    Args:
        input_file (str): Path to input CSV file
        output_file (str, optional): Path to output CSV file. If None, overwrite input file
        encoding (str, optional): Encoding of input file. Defaults to 'gb2312'
    """
    try:
        # Read the CSV file with specified encoding
        df = pd.read_csv(input_file, encoding=encoding)

        # Remove the last three rows
        if len(df) > 3:
            df = df.iloc[:-3]
        else:
            print(f"Warning: File {input_file} has {len(df)} rows or less, no rows were removed")

        # Determine output file name
        if output_file is None:
            output_file = input_file

        # Save with UTF-8 encoding
        df.to_csv(output_file, index=False, encoding='utf-8')
        print(f"Successfully processed {input_file} and saved to {output_file}")

    except Exception as e:
        print(f"Error processing file: {str(e)}", file=sys.stderr)
        sys.exit(1)


def main():
    """CLI entry point"""
    # Set up argument parser
    parser = argparse.ArgumentParser(
        description='Clean CSV files by removing the last three rows and converting to UTF-8')
    parser.add_argument('input_file', help='Input CSV file')
    parser.add_argument('output_file', nargs='?', default=None, help='Output CSV file (optional)')
    parser.add_argument('--encoding', default='gb2312', help='Input file encoding (default: gb2312)')

    # Parse arguments
    args = parser.parse_args()

    # Check if input file exists
    if not os.path.isfile(args.input_file):
        print(f"Error: Input file '{args.input_file}' does not exist", file=sys.stderr)
        sys.exit(1)

    # Process the CSV file
    process_csv(args.input_file, args.output_file, args.encoding)


if __name__ == "__main__":
    main()
