import os
from collections import defaultdict

def count_code_metrics(directory):
    total_lines = 0
    total_files = 0
    total_dirs = 0
    file_counts_by_extension = defaultdict(int)

    for dirpath, dirnames, filenames in os.walk(directory):
        total_dirs += 1
        for filename in filenames:
            if filename.endswith(('.go', '.sql', '.mod', '.sum', '.yml', '.yaml', 'Dockerfile', '.md')):
                total_files += 1
                extension = filename.split('.')[-1]
                file_counts_by_extension[extension] += 1
                
                with open(os.path.join(dirpath, filename), 'r', encoding='utf-8', errors='ignore') as file:
                    lines = file.readlines()
                    total_lines += len(lines)

    return total_dirs, total_files, file_counts_by_extension, total_lines

repo_path = '.'
dirs, files, file_counts, lines_of_code = count_code_metrics(repo_path)

print("RedHub-stats")
print(f'Total directories: {dirs}')
print(f'Total files: {files}')

print('File counts by extension:')
for ext, count in sorted(file_counts.items()):
    print(f'  .{ext}: {count}')

print(f'Total lines of code: {lines_of_code}')