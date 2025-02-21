#!/bin/bash

# Function to download SDF using CAS number
download_sdf() {
    CAS_ID=$1
    if [[ -z "$CAS_ID" ]]; then
        echo "Usage: $0 <CAS_ID>"
        exit 1
    fi

    # Convert CAS ID to PubChem CID
    CID=$(curl -s "https://pubchem.ncbi.nlm.nih.gov/rest/pug/compound/name/$CAS_ID/cids/TXT")
    
    if [[ -z "$CID" ]]; then
        echo "Error: No PubChem CID found for CAS ID $CAS_ID"
        exit 1
    fi

    echo "PubChem CID for CAS ID $CAS_ID: $CID"

    mkdir pubchem_downloads
    # Download SDF format file
    curl -s -o "./pubchem_downloads/${CAS_ID}.sdf" "https://pubchem.ncbi.nlm.nih.gov/rest/pug/compound/cid/${CID}/record/SDF?record_type=3d"

    if [[ -f "${CAS_ID}.sdf" ]]; then
        echo "Download successful: ${CAS_ID}.sdf"
    else
        echo "Error: Failed to download SDF file for CAS: $CAS_ID."
	echo $CAS_ID >> failed.txt
    fi
}

# Check if input file is provided
if [[ $# -ne 1 ]]; then
    echo "Usage: $0 <file_with_CAS_IDs>"
    exit 1
fi

FILE=$1

# Check if file exists
if [[ ! -f "$FILE" ]]; then
    echo "Error: File '$FILE' not found!"
    exit 1
fi

# Read CAS IDs line by line and process each
while IFS= read -r CAS_ID; do
    download_sdf "$CAS_ID"
done < "$FILE"

echo "🎉 Batch download complete!"
