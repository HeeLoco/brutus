FROM python:3.12.3

WORKDIR /ado

COPY src/server/ado/requirements.txt .
RUN pip install -r requirements.txt

EXPOSE 8000

COPY src/server/ado/ ./src/
CMD [ "python", "./src/main.py" ]