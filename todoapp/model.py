"""
model.py
"""
import client
import unicodedata


def connect_db(database):
    client.connect(database)

def save_db():
    return client.save()

def query(instruction):
    return client.query(instruction)

def show_tasks():
    return client.query("GET TASKS")

def new_task(title, task_id):   
    normalized_title = unicodedata.normalize("NFKD", title).encode("ascii", "ignore")

    stringquery = """addto tasks {"task_id":"%s", "title":"%s", "completed":"NIL", "show": "YES"}""" % (task_id, normalized_title)

    return client.query(stringquery)

def update_id(task_id):

    incr_task_id = int(task_id)+1
    update_id = """UPDATEINT TASK_ID %s""" % (incr_task_id)
    client.query(update_id)

def complete_task(task_id):

    stringquery = """NESTEDUPDATE tasks/%s/completed COMPLETE""" % (task_id)

    return client.query(stringquery)

def remove_task(task_id):

    stringquery = """NESTEDUPDATE tasks/%s/show NO""" % (task_id)

    return client.query(stringquery)