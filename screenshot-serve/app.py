from flask import Flask
from router import router as router_blueprint

app = Flask(__name__)
app.register_blueprint(router_blueprint)

if __name__ == '__main__':
    app.run()