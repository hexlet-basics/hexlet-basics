# frozen_string_literal: true

class DownloadExercise
  def self.run(lang_name)
    system("docker pull hexletbasics/exercises-#{lang_name}")
    system("rm -rf /var/tmp/hexletbasics/exercises-#{lang_name}")

    #FIXME can't pass data in volume inside docker.
    system("docker run --name exercises-#{lang_name} -v /var/tmp/hexletbasics/exercises-#{lang_name}:/out hexletbasics/exercises-#{lang_name}")
    system("docker cp exercises-#{lang_name}:/exercises-#{lang_name} /var/tmp/hexletbasics/")
    system("docker rm exercises-#{lang_name}")

    "/var/tmp/hexletbasics/exercises-#{lang_name}"
  end
end
