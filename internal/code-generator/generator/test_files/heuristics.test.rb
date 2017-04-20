    # Common heuristics
    ObjectiveCRegex = /^\s*(@(interface|class|protocol|property|end|synchronised|selector|implementation)\b|#import\s+.+\.h[">])/

    disambiguate ".asc" do |data|
      if /^(----[- ]BEGIN|ssh-(rsa|dss)) /.match(data)
        Language["Public Key"]
      elsif /^[=-]+(\s|\n)|{{[A-Za-z]/.match(data)
        Language["AsciiDoc"]
      elsif /^(\/\/.+|((import|export)\s+)?(function|int|float|char)\s+((room|repeatedly|on|game)_)?([A-Za-z]+[A-Za-z_0-9]+)\s*[;\(])/.match(data)
        Language["AGS Script"]
      end
    end

    fortran_rx = /^([c*][^abd-z]|      (subroutine|program|end|data)\s|\s*!)/i

    disambiguate ".f" do |data|
      if /^: /.match(data)
        Language["Forth"]
      elsif data.include?("flowop")
        Language["Filebench WML"]
      elsif fortran_rx.match(data)
        Language["FORTRAN"]
      end
    end

    disambiguate ".h" do |data|
      if ObjectiveCRegex.match(data)
        Language["Objective-C"]
      elsif (/^\s*#\s*include <(cstdint|string|vector|map|list|array|bitset|queue|stack|forward_list|unordered_map|unordered_set|(i|o|io)stream)>/.match(data) ||
        /^\s*template\s*</.match(data) || /^[ \t]*try/.match(data) || /^[ \t]*catch\s*\(/.match(data) || /^[ \t]*(class|(using[ \t]+)?namespace)\s+\w+/.match(data) || /^[ \t]*(private|public|protected):$/.match(data) || /std::\w+/.match(data))
        Language["C++"]
      end
    end

    disambiguate ".lsp", ".lisp" do |data|
      if /^\s*\((defun|in-package|defpackage) /i.match(data)
        Language["Common Lisp"]
      elsif /^\s*\(define /.match(data)
        Language["NewLisp"]
      end
    end

    disambiguate ".md" do |data|
      if /(^[-a-z0-9=#!\*\[|>])|<\//i.match(data) || data.empty?
        Language["Markdown"]
      elsif /^(;;|\(define_)/.match(data)
        Language["GCC machine description"]
      else
        Language["Markdown"]
      end
    end

    disambiguate ".ms" do |data|
      if /^[.'][a-z][a-z](\s|$)/i.match(data)
        Language["Groff"]
      elsif /(?<!\S)\.(include|globa?l)\s/.match(data) || /(?<!\/\*)(\A|\n)\s*\.[A-Za-z]/.match(data.gsub(/"([^\\"]|\\.)*"|'([^\\']|\\.)*'|\\\s*(?:--.*)?\n/, ""))
        Language["GAS"]
      else
        Language["MAXScript"]
      end
    end

    disambiguate ".mod" do |data|
      if data.include?('<!ENTITY ')
        Language["XML"]
      elsif /^\s*MODULE [\w\.]+;/i.match(data) || /^\s*END [\w\.]+;/i.match(data)
        Language["Modula-2"]
      else
        [Language["Linux Kernel Module"], Language["AMPL"]]
      end
    end

    disambiguate ".pro" do |data|
      if /^[^#]+:-/.match(data)
        Language["Prolog"]
      elsif data.include?("last_client=")
        Language["INI"]
      elsif data.include?("HEADERS") && data.include?("SOURCES")
        Language["QMake"]
      elsif /^\s*function[ \w,]+$/.match(data)
        Language["IDL"]
      end
    end

    disambiguate ".rpy" do |data|
      if /(^(import|from|class|def)\s)/m.match(data)
        Language["Python"]
      else
        Language["Ren'Py"]
      end
    end
