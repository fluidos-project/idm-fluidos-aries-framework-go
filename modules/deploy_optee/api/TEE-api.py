from flask import Flask, request, jsonify
import subprocess

app = Flask(__name__)

@app.route('/hello_world', methods=['GET'])
def run_hello_world():
    # Execute OP-TEE binary
    try:
        result = subprocess.run(["optee_example_hello_world"], capture_output=True, text=True)
        
        if result.returncode == 0:
            return jsonify({"message": "Success", "output": result.stdout}), 200
        else:
            return jsonify({"message": "Error executing optee_example_hello_world", "error": result.stderr}), 500

    except Exception as e:
        return jsonify({"message": "Exception occurred", "error": str(e)}), 500
        
        
@app.route('/dpabc_test', methods=['GET'])
def run_dpabc_test():
    # Execute dpabc_test command and returns the output
    try:
        result = subprocess.run(["dpabc_test"], capture_output=True, text=True, check=True)
        return jsonify({"message": "Success", "output": result.stdout})
    except subprocess.CalledProcessError as e:
        return jsonify({"message": "Error", "output": e.stdout, "error": e.stderr}), 500
    except FileNotFoundError:
        return jsonify({"message": "Error", "error": "dpabc_test command not found"}), 404
        
        
@app.route('/generate_key', methods=['POST'])
def generate_key():
    # Generate key
    data = request.json
    if "key_id" not in data:
        return jsonify({"error": "'key_id' required"}), 400

    key_id = data["key_id"]
    nattr = 5

    try:
        result = subprocess.run(
            ["/usr/bin/generate_key", key_id], capture_output=True, text=True
        )
        
        success = result.returncode == 0
        return jsonify({
            "success": success,
            "message": result.stdout.strip() if success else result.stderr.strip()
        }), 200 if success else 400

    except Exception as e:
        return jsonify({"error": "Error running generate_key", "message": str(e)}), 500
        
        
@app.route("/verify_signature", methods=["POST"])
def verify_signature():
    # Verify signature
    data = request.json
    required_keys = ["signature", "message", "public_key"]

    # Verify args
    if not all(key in data for key in required_keys):
        return jsonify({"error": "Invalid params"}), 400
        
    args = ["/usr/bin/verify_signature", data["signature"], data["message"], data["public_key"]]

    try:
        # Run script
        result = subprocess.run(args, capture_output=True, text=True)
        success = result.returncode == 0

        return jsonify({
            "success": success,
            "message": result.stdout.strip() if success else result.stderr.strip()
        }), 200 if success else 400
    except Exception as e:
        return jsonify({"error": "Error running verify_signature", "message": str(e)}), 500


@app.route("/generate_zktoken", methods=["POST"])
def generate_zktoken():
    # Generate ZkToken
    data = request.json
    required_keys = ["pkeyBase58", "indexes", "epoch", "nonce", "attributes"]

    # Verify args
    if not all(key in data for key in required_keys):
        return jsonify({"error": "Invalid params"}), 400

    args = ["/usr/bin/generate_zktoken", data["pkeyBase58"], data["indexes"], data["nonce"], data["epoch"]] + data["attributes"]

    try:
        # Run script
        result = subprocess.run(args, capture_output=True, text=True)
        success = result.returncode == 0

        return jsonify({
            "success": success,
            "message": result.stdout.strip() if success else result.stderr.strip()
        }), 200 if success else 400
        
    except Exception as e:
        return jsonify({"error": "Error running generate_zktoken", "message": str(e)}), 500
        

if __name__ == '__main__':
    app.run(host="10.0.2.15", port=5024)
