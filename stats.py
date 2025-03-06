import os

def count_lines_of_code(directory):
    total_lines = 0
    for dirpath, _, filenames in os.walk(directory):
        for filename in filenames:
            if filename.endswith(('.go', '.sql', '.mod', '.sum', '.yml', '.yaml', 'Dockerfile', '.md')):  # добавьте другие расширения по необходимости
                with open(os.path.join(dirpath, filename), 'r', encoding='utf-8', errors='ignore') as file:
                    lines = file.readlines()
                    total_lines += len(lines)
    return total_lines

repo_path = '.'
lines_of_code = count_lines_of_code(repo_path)
print(f'Total lines of code: {lines_of_code}')