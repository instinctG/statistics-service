.PHONY: up
up:
		echo "  --- running postgres on port :5432"
		docker-compose up --build
down:
		echo "  --- service is stopping"
		docker-compose down