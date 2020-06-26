# frozen_string_literal: true

class DownloadExercise
  def self.run(lang_name)
    system("docker pull hexletbasics/exercises-#{lang_name}")
    system("rm -rf tmp/hexletbasics/exercises-#{lang_name}")
    system("docker run --rm -v #{Dir.pwd}/tmp/hexletbasics/exercises-#{lang_name}:/out hexletbasics/exercises-#{lang_name} bash -c 'cp -r /exercises-#{lang_name}/* /out'")

    "tmp/hexletbasics/exercises-#{lang_name}"
  end
end
