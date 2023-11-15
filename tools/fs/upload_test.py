import requests
url = "http://localhost:8080/upload"
file_path = "README.md"
with open(file_path, "rb") as file:
    files = {'file': (file.name, file)}
    response = requests.post(url, files=files)
print(response.status_code)
print(response.text)