import requests
url = "http://127.0.0.1:8080/upload"
file_path = "README.md"
with open(file_path, "rb") as f:
    files = {"file": (file_path, f)}
    response = requests.post(url, files=files)
print(response.status_code)
print(response.text)