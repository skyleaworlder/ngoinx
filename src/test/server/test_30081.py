from flask import Flask, render_template_string

app = Flask(__name__, static_folder="./static_v3_test", static_url_path="/")

@app.route("/")
def server(**checkrst):
    return render_template_string("hahah this is 30081")

if __name__ == "__main__":
    app.run(host="127.0.0.1", port=30081, debug=True)