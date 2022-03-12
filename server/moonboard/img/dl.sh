#!/bin/bash
for i in {1..148}; do
	img=$i.png
	if [[ ! -f "$img" ]]; then
		if [[ $i -lt 41 || $i -gt 49 ]]; then
			curl -o $img https://moonboard.com/content/images/holds/h$i.png
		fi
	fi
done
for i in {1..80}; do
	img=w$i.png
	if [[ ! -f "$img" ]]; then
		if [[ $i -lt 10 ]]; then
			curl -o $img https://moonboard.com/content/images/holds/hw0$i.png
		else
			curl -o $img https://moonboard.com/content/images/holds/hw$i.png
		fi
	fi
done
