import csv
import json

f = open('movies.csv', 'r')
reader = csv.DictReader(f, delimiter="|", fieldnames = ( 
	"movie_id", "movie_title", "release_date", "video_release_date", "IMDb_URL", "unknown", "Action", "Adventure", "Animation", "Children", "Comedy", "Crime", "Documentary", "Drama", "Fantasy", "Film-Noir", "Horror", "Musical", "Mystery", "Romance", "Sci-Fi", "Thriller", "War", "Western" ) )
out = json.dumps( [ row for row in reader ] )

fo = open('movies.json', 'w')
fo.write(out)