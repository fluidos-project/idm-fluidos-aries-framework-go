import requests
import json

# API base URL
BASE_URL = "http://localhost:3000"

def enroll_admin():
    url = f"{BASE_URL}/enroll-admin"
    response = requests.post(url)
    if response.status_code == 200:
        print("Admin enrolled successfully")
        return True
    else:
        print(f"Error enrolling admin: {response.status_code}")
        return False

def connect_network():
    url = f"{BASE_URL}/connect-network"
    response = requests.post(url)
    if response.status_code == 200:
        print("Connected to network successfully")
        return True
    else:
        print(f"Error connecting to network: {response.status_code}")
        return False

def read_from_dht(key):
    url = f"{BASE_URL}/dht/{key}"
    response = requests.get(url)
    if response.status_code == 200:
        return response.json()
    else:
        print(f"Error reading from DHT: {response.status_code}")
        return None

def write_to_dht(key, value):
    url = f"{BASE_URL}/dht"
    data = {"key": key, "value": value}
    headers = {"Content-Type": "application/json"}
    response = requests.post(url, data=json.dumps(data), headers=headers)
    if response.status_code == 201:
        return response.json()
    else:
        print(f"Error writing to DHT: {response.status_code}")
        return None

def get_admin_identity():
    response = requests.get(f"{BASE_URL}/admin-identity")
    if response.status_code == 200:
        return response.json()['adminIdentity']
    else:
        print(f"Error getting admin identity: {response.status_code}")
        return None

# Example usage
if __name__ == "__main__":
    # Enroll admin
    enroll_admin()


    # Get admin identity
    admin_identity = get_admin_identity()
    if admin_identity:
        print("Admin identity retrieved successfully")
        print(json.dumps(admin_identity, indent=2))
    else:
        print("Failed to retrieve admin identity")


    # Connect to network
    connect_network()

    # Write to DHT
    write_result = write_to_dht("testKey", "testValue")
    if write_result:
        print("Write successful:", write_result)

    # Read from DHT
    read_result = read_from_dht("testKey")
    if read_result:
        print("Read successful:", read_result)

    