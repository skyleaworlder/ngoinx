from flask import Flask, render_template_string

app = Flask(__name__, static_folder="./static_v1_test", static_url_path="/")

@app.route("/api/v1/test")
def server(**checkrst):
    return render_template_string("hahah this is 10081")

if __name__ == "__main__":
    app.run(host="127.0.0.1", port=10081, debug=True)