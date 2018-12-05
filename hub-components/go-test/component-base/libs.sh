
echoerr() { echo "$@" 1>&2; }

__var_prefix="_WORKFLOW_"
hub_var() {
    local __var_name="$1"
    local __default_var_name="$__var_prefix$__var_name"
    local __required="$2"
    local __result_var="$3"
    local __result_value=

    if [[ "${!__var_name}" != "" ]];then
        __result_value="${!__var_name}"
    elif [[ "${!__default_var_name}" != "" ]]; then
        __result_value="${!__default_var_name}"
    elif [[ "$__required" == "true" ]]; then
        echoerr "required variable '$__var_name' not exist!"
        exit 1
    fi

    if [[ "$__result_var" ]]; then
        eval $__result_var="'$__result_value'"
    else
        echo "$__result_value"
    fi
}

hub_default_var() {
    local __default_value="$2"
    local __result_var="$3"
    local __result_value=
    if [[ -z "$__result_var" ]]; then
        __result_value=`hub_var "$1"`
        if [[ "$__result_value" ]]; then
            echo "$__result_value"
        else
            echo "$__default_value"
        fi
    else
        hub_var "$1" "false" "$__result_var"
        if [[ -z "${!__result_var}" ]]; then
            eval $__result_var="'$__default_value'"
        fi
    fi
}

hub_git_clone() {
  GIT_REPO_URL=`hub_var GIT_CLONE_URL true`
  GIT_REPO_REF=`hub_default_var GIT_REF master`
  GIT_REPO_DIR=`hub_var GIT_REPO_DIR true`

  GIT_SSH_COMMAND="ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no" git clone --recurse-submodules $GIT_REPO_URL "$GIT_REPO_DIR"

  cd $GIT_REPO_DIR
  git checkout "$GIT_REPO_REF" --
}
