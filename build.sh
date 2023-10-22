if [ ! -f .env.production ]; then
    echo ".env.production file not found!"
    exit 1
fi

while IFS= read -r line; do
    if [[ ! $line =~ ^# && $line =~ = ]]; then
		echo "exporting $line"
        export "$line"
    fi
done < .env.production

goreleaser release