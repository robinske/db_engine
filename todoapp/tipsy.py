"""
tipsy.py -- A flask-based todo list
"""
from flask import Flask, render_template, redirect, request
import model
import json

app = Flask(__name__)

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
    description = "You can update this task's description"    
    
    task_id = model.query("GET TASK_ID")
    model.update_id(task_id)
    task_id_updated = model.query("GET TASK_ID")

    model.new_task(title, description, task_id_updated)

    return redirect("/")

if __name__ == "__main__":
    model.connect_db("stage1/tasks.db")
    app.run(debug=True)