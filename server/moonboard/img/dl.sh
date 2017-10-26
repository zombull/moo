#!/bin/bash
for i in {1..149}; do
	if [[ $i -lt 41 || $i -gt 49 ]]; then
		curl -O https://moonboard.com/wp-content/plugins/moonboard/inc/img/$i.png
	fi
done
