"""Setup script for TodoList CLI"""

from setuptools import setup, find_packages

with open("README.md", "r", encoding="utf-8") as fh:
    long_description = fh.read()

setup(
    name="todolist-cli",
    version="1.0.0",
    author="Your Name",
    description="A smart command-line todo list manager with AI Agent",
    long_description=long_description,
    long_description_content_type="text/markdown",
    packages=find_packages(),
    classifiers=[
        "Development Status :: 4 - Beta",
        "Intended Audience :: Developers",
        "Topic :: Software Development :: Build Tools",
        "License :: OSI Approved :: MIT License",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.8",
        "Programming Language :: Python :: 3.9",
        "Programming Language :: Python :: 3.10",
        "Programming Language :: Python :: 3.11",
    ],
    python_requires=">=3.8",
    install_requires=[
        "click>=8.1.7",
        "anthropic>=0.39.0",
        "python-dotenv>=1.0.0",
        "rich>=13.7.0",
    ],
    entry_points={
        "console_scripts": [
            "todo=src.main:main",
        ],
    },
)
