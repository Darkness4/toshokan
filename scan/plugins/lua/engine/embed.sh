#!/bin/bash
# Embed lua files into C strings to avoid any security issues.

convert_to_header() {
    input_file="$1"
    output_header_file="$2"

    # Generate a valid variable name from the input file name
    var_name=$(basename "$input_file" | sed 's/[^a-zA-Z0-9_]/_/g')_data

    # Start writing the header file
    {
        echo "#ifndef ${var_name}_H"
        echo "#define ${var_name}_H"
        echo ""
        echo "const char ${var_name}[] = \\"

        # Read the input file line by line
        while IFS= read -r line; do
            # Convert each line to a C string literal
            echo -n "    \""
            # Escape special characters for C string literals
            echo "$line" | sed 's/\\/\\\\/g; s/"/\\"/g; s/$/\\n"/'
        done < "$input_file"

        # Close the C string literal and header guard
        echo ";"
        echo ""
        echo "#endif // ${var_name}_H"
    } > "$output_header_file"

    echo "Embedded $input_file into $output_header_file"
}

SCRIPT_DIR=$(dirname "$0")

convert_to_header "$SCRIPT_DIR/msgpack.lua" "$SCRIPT_DIR/msgpack.lua.h"
