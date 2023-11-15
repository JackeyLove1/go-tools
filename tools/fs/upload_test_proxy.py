import requests

url = "http://localhost:6789/proxy/upload"
file_path = "README.md"
with open(file_path, "rb") as file:
    headers = {"address": "http://localhost:8080/upload"}
    files = {'file': (file.name, file)}
    response = requests.post(url, files=files, headers=headers)
print(response.status_code)
print(response.text)
