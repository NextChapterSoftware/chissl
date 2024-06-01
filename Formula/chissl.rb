# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Chissl < Formula
  desc "HTTPS reverse tunnel server/client"
  homepage "https://github.com/NextChapterSoftware/chissl"
  version "1.2"
  license "MIT"

  on_macos do
    url "https://github.com/NextChapterSoftware/chissl/releases/download/v1.2/chissl_Darwin_all.zip"
    sha256 "dea32d02ecfc0fad844b86bf36596f495f41f308ca2ada61861c44a4c6bfc1ad"

    def install
      bin.install "chissl"
    end
  end

  on_linux do
    on_intel do
      if Hardware::CPU.is_64_bit?
        url "https://github.com/NextChapterSoftware/chissl/releases/download/v1.2/chissl_Linux_x86_64.zip"
        sha256 "feae7e4eec1a3192c6ebc2d58f9b6972f86a79578ae04ef4160b14fd0c0ef950"

        def install
          bin.install "chissl"
        end
      end
    end
    on_arm do
      if !Hardware::CPU.is_64_bit?
        url "https://github.com/NextChapterSoftware/chissl/releases/download/v1.2/chissl_Linux_armv6.zip"
        sha256 "f596ce667d951ce6cb3e45a780871a59b26809fcadcced15e08f43c6330685e8"

        def install
          bin.install "chissl"
        end
      end
    end
    on_arm do
      if Hardware::CPU.is_64_bit?
        url "https://github.com/NextChapterSoftware/chissl/releases/download/v1.2/chissl_Linux_arm64.zip"
        sha256 "577f90272ad7cfad15f6e9ea8e679bad8921223ef2b6bd0f815bea6464f72d3a"

        def install
          bin.install "chissl"
        end
      end
    end
  end
end