class Create(self):

	collection = None
	# get this from the user/person updating the db

	def __init__(self, collection):
		# throw this into the database being used
		# create a collection with name obj_name
		# give it characteristics (like unique ID keys, predetermined columnn names, flexibility to add more column names)
		# how does the user specify this? "raw input"?
		# am i writing to a file or writing to an interpreter that writes to the db? i.e. so that it can be inputted in binary...


class Read(self):
	
	def __init__(self, db): #collection too??
		self.db = db

	def open(db):
		# open and read the file...
		# how to do this efficiently without looking through the whole thing
		pass

	def filtered_read_method(given_these_parameters):
		pass
		# one possible method

class Insert(self):

	def __init__(self):
		pass

	def insert_new_record(self, collection):
		# must pass in the name of the collection/table/object you want to put the record/row into
		pass

class Select(self):

	def __init__(self, collection):
		# get a record back out of the collection in a given db
		pass



