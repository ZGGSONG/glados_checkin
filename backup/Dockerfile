FROM python:3.9-alpine

WORKDIR /app

ADD ./requirements.txt /app/.
COPY *.py /app/.

RUN pip3 install -r requirements.txt

CMD ["python3", "main.py"]