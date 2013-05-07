"""
tipsy.py -- A flask-based todo list
"""
from flask import Flask, render_template, redirect, request
import model
import json

app = Flask(__name__)

# @app.before_request
# def before_request():
#     model.connect_db("todoapp/tasks.db")

@app.route("/", methods = ["GET"])
def index():
    
    json_string = model.show_tasks()
    tasks = json.loads(json_string)
    return render_template("index.html", tasks=tasks)

@app.route("/complete_task", methods= ["POST"])
def complete_task():

    task_id = str(request.form.get("task_id"))
    model.complete_task(task_id)
    
    return "Task marked as complete"

@app.route("/remove_task", methods= ["POST"])
def remove_task():

    task_id = str(request.form.get("task_id"))
    model.remove_task(task_id)
    
    return "Task will be removed"

@app.route("/new_task", methods = ["POST"])
def add_task():
    
    title = request.form["new_task_title"]  
    
    task_id = model.query("GET TASK_ID")
    model.update_id(task_id)
    task_id_updated = int(task_id) + 1

    model.new_task(title, task_id_updated)

    return redirect("/")

if __name__ == "__main__":
    model.connect_db("todoapp/tasks.db") # this is causing it to load twice...messing with the rest of it
    app.run(debug=True)
