#!/bin/bash

# PubChem base URL
BASE_URL="https://pubchem.ncbi.nlm.nih.gov/rest/pug/compound/cid"

# Check if input file is provided
if [ -z "$1" ]; then
    echo "Usage: $0 <cid_list.txt>"
    echo "CID is downloaded in SDF format"
    exit 1
fi

CID_FILE="$1"

# Create an output directory
OUTPUT_DIR="pubchem_downloads"
mkdir -p "$OUTPUT_DIR"

# Read CIDs from the file and fetch data
while IFS= read -r CID; do
    if [[ -n "$CID" ]]; then
        echo "Downloading CID: $CID in $OUTPUT_FORMAT format..."
	OUTPUT_FILE="${OUTPUT_DIR}/compound_${CID}.sdf"
        URL="${BASE_URL}/${CID}/record/SDF?record_type=3d"
        
	# Download the data
        curl -s "$URL" -o "$OUTPUT_FILE"
        
        if [[ -s "$OUTPUT_FILE" ]]; then
            echo "Saved: $OUTPUT_FILE"
        else
            echo "Error: Could not fetch data for CID $CID"
            rm -f "$OUTPUT_FILE" # Remove empty files
        fi
    fi
done < "$CID_FILE"

echo "Download complete. Files saved in '$OUTPUT_DIR'."
