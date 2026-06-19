import re
import sys
import os

def main():
    if len(sys.argv) < 3:
        sys.exit(1)

    changelog_path = sys.argv[1]
    tag_name = sys.argv[2]
    
    version = tag_name.lstrip('v')
    
    with open(changelog_path, 'r', encoding='utf-8') as f:
        content = f.read()
        
    lines = content.splitlines()
    header_idx = -1
    pattern = re.compile(rf'^##\s+\[{re.escape(version)}\]')
    
    for idx, line in enumerate(lines):
        if pattern.match(line):
            header_idx = idx
            break
            
    if header_idx == -1:
        sys.exit(1)
        
    title = lines[header_idx].replace('##', '').strip()
    
    body_lines = []
    for line in lines[header_idx + 1:]:
        if line.startswith('## ['):
            break
        body_lines.append(line)
        
    body = '\n'.join(body_lines).strip()
    
    if 'GITHUB_OUTPUT' in os.environ:
        with open(os.environ['GITHUB_OUTPUT'], 'a', encoding='utf-8') as fh:
            fh.write(f"title={title}\n")
            
    with open('release_body.md', 'w', encoding='utf-8') as fh:
        fh.write(body)
        
if __name__ == '__main__':
    main()
