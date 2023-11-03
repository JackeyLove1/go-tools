import requests
url = "http://nas.admin.byted.org/proxy/http/n250-076-160:8080/upload"
# url = "http://localhost:8080/upload"
file_path = "README.md"
with open(file_path, "rb") as file:
    files = {'file': (file.name, file)}
    response = requests.post(url, files=files)
print(response.status_code)
print(response.text)