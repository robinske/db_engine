import csv
import json

f = open('sample.csv', 'r')
reader = csv.DictReader(f, fieldnames = ( "id","name","lat","lng" ) )
out = json.dumps( [ row for row in reader ] )
print out