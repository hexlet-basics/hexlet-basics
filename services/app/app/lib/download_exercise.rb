# frozen_string_literal: true

class DownloadExercise
  def self.run(lang_name)
    system("docker pull hexletbasics/exercises-#{lang_name}")
    system("rm -rf /var/tmp/hexletbasics/exercises-#{lang_name}")
    system("docker run --rm -v /var/tmp/hexletbasics/exercises-#{lang_name}:/out hexletbasics/exercises-#{lang_name} bash -c 'cp -r /exercises-#{lang_name}/* /out'")

    "/var/tmp/hexletbasics/exercises-#{lang_name}"
  end
end
