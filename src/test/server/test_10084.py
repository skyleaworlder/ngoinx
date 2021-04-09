from flask import Flask, render_template_string

app = Flask(__name__, static_folder="./static_v1_food", static_url_path="/")

@app.route("/")
def server(**checkrst):
    return render_template_string("hahah this is 10084")

if __name__ == "__main__":
    app.run(host="127.0.0.1", port=10084, debug=True)