# frozen_string_literal: true

class DockerImageExerciseLoader
  def run(lang_name)
    system("docker pull hexletbasics/exercises-#{lang_name}")
    system("rm -rf tmp/hexletbasics/exercises-#{lang_name}")
    system("docker run --rm -v #{Dir.pwd}/tmp/hexletbasics/exercises-#{lang_name}:/out hexletbasics/exercises-#{lang_name} bash -c 'cp -r /exercises-#{lang_name}/* /out'")
  end
end
