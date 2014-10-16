package ro.redeul.google.go.lang.psi.resolve;

import com.intellij.psi.PsiElement;

/**
* Created by mihai on 10/14/14.
*/
interface GoReferenceSolver<R extends GoReferenceWithSolver<?, S, R>, S extends GoReferenceSolver<R, S>> {

    public PsiElement resolve(R reference);

    public Object[] getVariants(R reference);
}
