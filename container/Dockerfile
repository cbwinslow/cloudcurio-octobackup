FROM nvidia/cuda:12.2.2-runtime-ubuntu22.04
RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y git curl ca-certificates python3 python3-pip python3-venv && rm -rf /var/lib/apt/lists/*
WORKDIR /opt/review
COPY requirements.txt .
RUN pip3 install --no-cache-dir -r requirements.txt
COPY review_runner.py /opt/review/review_runner.py
ENV PYTHONUNBUFFERED=1 PIP_DISABLE_PIP_VERSION_CHECK=1 PYTHONDONTWRITEBYTECODE=1
WORKDIR /work
ENTRYPOINT ["python3","/opt/review/review_runner.py"]
