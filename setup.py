from setuptools import setup, find_packages

setup(
    name="adk-router",
    version="0.1.0",
    description="Provider-agnostic ADK extension for OpenAI-compatible LLM providers",
    author="ADK Router Team",
    packages=find_packages(),
    install_requires=[
        "google-adk>=0.1.0",
        "litellm>=1.0.0",
        "python-dotenv>=1.0.0",
    ],
    python_requires=">=3.9",
    license="Apache-2.0",
    classifiers=[
        "Development Status :: 3 - Alpha",
        "Intended Audience :: Developers",
        "License :: OSI Approved :: Apache Software License",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.9",
        "Programming Language :: Python :: 3.10",
        "Programming Language :: Python :: 3.11",
        "Programming Language :: Python :: 3.12",
    ],
)