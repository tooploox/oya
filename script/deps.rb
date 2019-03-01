require 'tmpdir'

def include?(import_path)
            import_path.start_with?("github.com/tooploox/oya/pkg") &&
              !import_path.end_with?("/internal") &&
              !import_path.end_with?("/fixtures") &&
              !import_path.end_with?("/debug")
end

deps = `go list -f '{{.ImportPath}} -> {{ .Imports }}' github.com/bilus/...`

line_rx = /^([^\s]+)\s+->\s+\[([^\]]*)\]\s*$/

Dir.mktmpdir do |dir|
  output = File.open(File.join(dir, "output"), 'w+')
  output.puts """
@startuml
digraph G {
"""

  deps.split("\n").flat_map do |line|
    matches = line.match(line_rx)
    package = matches[1]
    if include?(package)
      imports = matches[2].split(" "). \
                  select { |import_path| include?(import_path) }. \
                  map { |import_path| "\"#{import_path}\"" }. \
                  join(", ")
      if !imports.empty?
        output.puts "  \"#{package}\" -> #{imports}"
      end
    end
  end

  output.puts """
}
@enduml
"""

  output.flush

  image = File.open(File.join(dir, "image.svg"), 'w+')

  `cat #{output.path} | docker run --rm -i think/plantuml > #{image.path}`
  `open #{image.path}`
end
