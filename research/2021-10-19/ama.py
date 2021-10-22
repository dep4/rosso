import requests

headers = {'User-Agent': 'Mozilla'}
resp = requests.get('https://www.amazon.com/dp/B07K5214NZ', headers=headers)
print(resp.text)
